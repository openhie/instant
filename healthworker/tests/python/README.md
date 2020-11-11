# Tests

> This folder of sample tests will eventually be revised or removed. It remains for archival purposes'.

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

The test data was taken from the iHRIS repository: https://github.com/iHRIS/iHRIS/tree/master/resources/demo

```
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-Country.json --output Country.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-District.json --output District.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-Facility.json --output Facility.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-Position.json --output Position.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-Practitioner-HBB.json --output Practitioner-HBB.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-Practitioner.json --output Practitioner.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-PractitionerRole.json --output PractitionerRole.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-Region.json --output Region.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/bundle-Salary.json --output Salary.json
curl https://raw.githubusercontent.com/iHRIS/iHRIS/master/resources/demo/staff.iHRISRelationship.json --output staffiHRISRelationship.json
```

Launch your FHIR server of choice. Upload the records to a FHIR server edit the `postbundles.sh` file to use your path to the data, e.g.:
```sh
for i in $( ls data); do
    echo $i
    curl http://localhost:8080/baseR4 --data-binary "@/Users/richard/src/github.com/openhie/tests/examples/components/healthworker/testdata/$i" -H "Content-Type: application/fhir+json"
done
```

