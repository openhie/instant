# Tests for health worker registry component

These is a sample of how to construct `.feature` files, the associated steps in Python to run them, and how to load sample data.

# Feature files preparation

`.feature` files are compatible representations of the Gherkin language for behavior-driven development. They may be written by analysts who define features in scenarios. 

```
Feature: Basic Health Worker Registry Queries

  Scenario: Query for HW records
     Given fhir server returns a capability statement for practitioner
      When practitioner resources exist on the fhir server
      Then get one practitioner resource
```

The `given`, `when`, and `then` statements become function decorators in the associated steps file which define the actual code to run. 

# Steps in code

Most languages have a library to support feature files. This example uses the behave package for Python. To use it, install Python and behave

Modify steps/required.py to your preferred FHIR server using an environment variable `FHIR_SERVER` or accept the default in the file.

Then run `behave`.
```txt
$ behave
# output:
Feature: Basic Health Worker Registry Queries # required.feature:1

  Scenario: Query for HW records                                      # required.feature:3
    Given fhir server returns a capability statement for practitioner # steps/required.py:9 0.114s
    When practitioner resources exist on the fhir server              # steps/required.py:21 0.074s
    Then get one practitioner resource                                # steps/required.py:33 0.067s

1 feature passed, 0 failed, 0 skipped
1 scenario passed, 0 failed, 0 skipped
3 steps passed, 0 failed, 0 skipped, 0 undefined
Took 0m0.255s
```

# Example of data preparation

This process uses [Synthea](https://github.com/synthetichealth/synthea) to generate fake data for testing workflows and components.

Ensure Java and gradle are installed. Clone and check the Synthea repository.

```sh
git clone https://github.com/synthetichealth/synthea.git
cd synthea
./gradlew build check test
```

Create 10 records with seed 1234 for reproducibility.
```sh
./run_synthea -s 1234 -p 10
```

Copy the records to the `data` folder.
```sh
cp output/fhir/*.json ~/src/github.com/openhie/tests/examples/components/healthworker/data/
```

Launch your FHIR server of choice. Upload the records to a FHIR server.
```sh
for i in $( ls data); do
    echo $i
    curl http://localhost:8080/baseR4 --data-binary "@/Users/richard/src/github.com/openhie/tests/examples/components/healthworker/data/$i" -H "Content-Type: application/fhir+json"
done
```

