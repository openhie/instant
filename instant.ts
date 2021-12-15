'use strict'

import * as commandLineArgs from 'command-line-args'
import * as glob from 'glob'
import * as fs from 'fs'
import * as child from 'child_process'
import * as util from 'util'
import * as path from 'path'
import { env } from 'process'

const exec = util.promisify(child.exec)

interface PackageInfo {
  metadata: {
    id: string
    name: string
    description: string
    version: string
    dependencies: string[]
    environmentVariables: object
  }
  path: string
}

interface PackagesMap {
  [packageID: string]: PackageInfo
}

function getInstantOHIEPackages(): PackagesMap {
  const packages: PackagesMap = {}
  let pathRegex = 'instant.json'
  let paths = []
  let nestingLevel = 0

  while (nestingLevel < 5) {
    pathRegex = "*/" + pathRegex
    const nestedPackages = glob.sync(pathRegex)
    
    paths = paths.concat(nestedPackages)
    nestingLevel += 1
  }

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
  function resolveDeps(id, currentStack) {
    if (currentStack.includes(id)) throw Error(`Circular dependency present for id ${id}`)
    currentStack.push(id)

    if (allPackages[id] && allPackages[id].metadata) {
      if (
        !allPackages[id].metadata.dependencies ||
        !allPackages[id].metadata.dependencies.length
      ) return [id]
    } else {
      throw Error(`Package ${id} does not exist or the metadata is invalid`)
    }

    const orderedIds = []
    const currentStackClone = currentStack.slice()

    allPackages[id].metadata.dependencies.forEach(dependency => {
      const ids = resolveDeps(dependency, currentStackClone)
      orderedIds.push(...ids)
    })
    orderedIds.push(id)
    return orderedIds
  }

  let orderedPackageIds = []
  chosenPackageIds.forEach(packageId => {
    let packageIds = orderedPackageIds.concat(resolveDeps(packageId, []))
    orderedPackageIds = packageIds.filter((id, index) => packageIds.indexOf(id) == index)
  })
  return orderedPackageIds
}

const logPackageDetails = (packageInfo: PackageInfo) => {
  console.log(`------------------------------------------------------------\nConfig Details: ${packageInfo.metadata.name} (${packageInfo.metadata.id})\n------------------------------------------------------------`)
  const envVars = []
  for(let envVar in packageInfo.metadata.environmentVariables) {
    envVars.push({"Environment Variable": envVar, "Default Value": packageInfo.metadata.environmentVariables[envVar], "Updated Value": env[envVar]})
  }
  console.table(envVars)
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
        },
        {
          name: 'only',
          alias: 'o',
          type: Boolean
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
      throw new Error(`Deploy - Unknown package id in list: ${chosenPackageIds}`)
    }

    if (chosenPackageIds.length < 1) {
      chosenPackageIds = Object.keys(allPackages)
    }

    if (!mainOptions.only) {
      // Order the packages such that the dependencies are instantiated first
      chosenPackageIds = orderPackageIds(allPackages, chosenPackageIds)
    }

    if (['destroy', 'down'].includes(main.command)) {
      chosenPackageIds.reverse()
    }

    console.log(
      `Selected package IDs to operate on: ${chosenPackageIds.join(', ')}`
    )

    switch (mainOptions.target) {
      case 'docker':
        for (const id of chosenPackageIds) {
          logPackageDetails(allPackages[id])
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
      throw new Error(`Testing - Unknown package id in list: ${chosenPackageIds}`)
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
