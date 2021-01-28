'use strict'

import * as commandLineArgs from 'command-line-args'
import * as glob from 'glob'
import * as fs from 'fs'
import * as child from 'child_process'
import * as util from 'util'
import * as path from 'path'

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

async function runTests(path: string) {
  const cmd = `node_modules/.bin/cucumber-js ${path}`

  try {
    const promise = exec(cmd)
    if (promise.child) {
      promise.child.stdout.on('data', (data) => console.log(data))
      promise.child.stderr.on('data', (data) => console.error(data))
    }
    await promise
  } catch (err) {
    console.error(`Error: Tests at ${path} returned an error`)
    console.log(err.stdout)
    console.log(err.stderr)
  }
}

const orderPackageIds = (allPackages, chosenPackageIds) => {
  const orderedPackageIds = []
  const packagesWithDependencies = []

  function collectPackagesWithDependencies(chosenIds) {
    let inexplicitDependencies = []

    chosenIds.forEach(id => {
      if (id && allPackages[id] && allPackages[id].metadata) {
        if (
          !allPackages[id].metadata.dependencies ||
          !allPackages[id].metadata.dependencies.length
        ) {
          if (id === 'core') {
            orderedPackageIds.unshift(id)
          } else {
            orderedPackageIds.push(id)
          }
        } else {
          allPackages[id].metadata.dependencies.forEach(dependency => {
            if (!Object.keys(allPackages).includes(dependency)) {
              throw Error(`Dependency ${dependency} for package ${id} does not exist`)
            }
            if (!chosenIds.includes(dependency)) {
              inexplicitDependencies.push(dependency)
            }
          })
          packagesWithDependencies.push(id)
        }
      } else {
        throw Error(`Package ${id} does not exist or the metadata is invalid`)
      }
    })

    if (inexplicitDependencies.length) {
      collectPackagesWithDependencies(inexplicitDependencies)
    } else {
      return
    }
  }

  collectPackagesWithDependencies(chosenPackageIds)

  while (packagesWithDependencies.length) {
    const currentPackagesLength = packagesWithDependencies.length

    for (let index = 0; index < packagesWithDependencies.length; index++) {
      const id = packagesWithDependencies[index]
      let containDependencies = true

      allPackages[id].metadata.dependencies.forEach(dependency => {
        if (!orderedPackageIds.includes(dependency)) {
          containDependencies = false
        }
      })

      if (containDependencies) {
        orderedPackageIds.push(id)
        packagesWithDependencies.splice(index, 1)
      }
    }
    /*
      If circular dependencies are present, the array of packages with dependencies will
      not change after the loop above finishes
    */
    if (currentPackagesLength == packagesWithDependencies.length) {
      throw Error('Error! Circular dependencies present')
    }
  }
  return orderedPackageIds
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

    // Order the packages such that the dependencies are instantiated first
    chosenPackageIds = orderPackageIds(allPackages, chosenPackageIds)

    if (['destroy', 'down'].includes(main.command)) {
      chosenPackageIds.reverse()
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

    // Order the packages such that the dependencies are instantiated first
    chosenPackageIds = orderPackageIds(allPackages, chosenPackageIds)

    console.log(`Running tests for packages: ${chosenPackageIds.join(', ')}`)
    console.log(`Using host: ${testOptions.host}:${testOptions.port}`)

    for (const id of chosenPackageIds) {
      const features = path.resolve(allPackages[id].path, 'features')
      await runTests(features)
    }
  }
})()
