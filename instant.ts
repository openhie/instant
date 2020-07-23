'use strict'

import * as commandLineArgs from 'command-line-args'
import * as glob from 'glob'
import * as fs from 'fs'
import * as child from 'child_process'
import * as util from 'util'

const exec = util.promisify(child.exec)

interface PackageInfo {
  metadata: {
    id: string
    name: string
    description: string
    version: string
    dependencies: string[]
  }
  path: string
}

interface PackagesMap {
  [packageID: string]: PackageInfo
}

function getInstantOHIEPackages(): PackagesMap {
  const packages: PackagesMap = {}
  const paths = glob.sync('*/instant.json')

  for (const path of paths) {
    const metadata = JSON.parse(fs.readFileSync(path).toString())
    packages[metadata.id] = {
      metadata,
      path: path.replace('instant.json', '')
    }
  }

  return packages
}

async function runBashScript(path: string, filename: string, args: string[]) {
  const cmd = `bash ${path}${filename} ${args.join(' ')}`
  console.log(`Executing: ${cmd}`)

  try {
    const promise = exec(cmd)
    if (promise.child) {
      promise.child.stdout.on('data', (data) => console.log(data))
      promise.child.stderr.on('data', (data) => console.error(data))
    }
    await promise
  } catch (err) {
    console.error(`Error: Script ${filename} returned an error`)
    console.log(err.stdout)
    console.log(err.stderr)
  }
}

// Main script execution
;(async () => {
  const allPackages = getInstantOHIEPackages()
  console.log(
    `Found ${Object.keys(allPackages).length} packages: ${Object.values(
      allPackages
    )
      .map((p) => p.metadata.id)
      .join(', ')}`
  )

  const main = commandLineArgs(
    [
      {
        name: 'command',
        defaultOption: true
      }
    ],
    {
      stopAtFirstUnknown: true
    }
  )

  let argv = main._unknown || []

  // main commands
  if (['init', 'up', 'down', 'destroy'].includes(main.command)) {
    const mainOptions = commandLineArgs(
      [
        {
          name: 'target',
          alias: 't',
          defaultValue: 'docker'
        }
      ],
      { argv, stopAtFirstUnknown: true }
    )

    console.log(`Target environment is: ${mainOptions.target}`)

    argv = mainOptions._unknown || []
    let chosenPackageIds = argv

    if (
      !chosenPackageIds.every((id) => Object.keys(allPackages).includes(id))
    ) {
      throw new Error('Unknown package id')
    }

    if (chosenPackageIds.length < 1) {
      chosenPackageIds = Object.keys(allPackages)
    }

    console.log(
      `Selected package IDs to operate on: ${chosenPackageIds.join(', ')}`
    )

    switch (mainOptions.target) {
      case 'docker':
        for (const id of chosenPackageIds) {
          await runBashScript(`${allPackages[id].path}docker/`, 'compose.sh', [
            main.command
          ])
        }
        break
      case 'k8s':
      case 'kubernetes':
        for (const id of chosenPackageIds) {
          await runBashScript(
            `${allPackages[id].path}kubernetes/main/`,
            'k8s.sh',
            [main.command]
          )
        }
        break
      default:
        throw new Error("Unknown value given for option 'target'")
    }
  }

  // test command
  if (main.command === 'test') {
    const testOptions = commandLineArgs(
      [
        {
          name: 'host',
          alias: 'h',
          defaultValue: 'localhost'
        },
        {
          name: 'port',
          alias: 'p',
          defaultValue: '5000'
        }
      ],
      { argv, stopAtFirstUnknown: true }
    )

    argv = testOptions._unknown || []
    let chosenPackageIds = argv

    if (
      !chosenPackageIds.every((id) => Object.keys(allPackages).includes(id))
    ) {
      throw new Error('Unknown package id')
    }

    if (chosenPackageIds.length < 1) {
      chosenPackageIds = Object.keys(allPackages)
    }

    console.log(`Running tests for packages: ${chosenPackageIds.join(', ')}`)
    console.log(`Using host: ${testOptions.host}:${testOptions.port}`)

    for (const id of chosenPackageIds) {
      await runBashScript(`${allPackages[id].path}`, 'test.sh', [
        `${testOptions.host}:${testOptions.port}`
      ])
    }
  }
})()
