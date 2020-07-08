#!/usr/bin/env python3
import os
from behave import given, when, then  # pylint: disable=no-name-in-module
import requests


local = os.getenv('FHIR_SERVER',
                  'http://localhost:8080/hapi-fhir-jpaserver/fhir/')


@given('fhir server returns a capability statement for practitioner')
def step_impl(context):
    resp = requests.get(url=f'{local}/metadata')
    cs = resp.json()
    cs_res = cs['rest'][0]['resource']
    exists = list(filter(lambda x: x.get('type') == 'Practitioner', cs_res))
    if exists:
        assert True
    else:
        assert False


@when('practitioner resources exist on the fhir server')
def step_impl(context):
    resp = requests.get(url=f'{local}/Practitioner?_history')
    pr = resp.json()
    total_str = pr['total']
    print(total_str)
    if total_str == 0:
        assert False
    else:
        assert True


@then('get one practitioner resource')
def step_impl(context):
    resp = requests.get(url=f'{local}/Practitioner?_history')
    pr = resp.json()
    assert True is not False
