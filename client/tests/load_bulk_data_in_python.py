#!/usr/bin/env python
# coding: utf-8

# # Load bulk data in Python
# 
# This example is available in the Jupyter notebook at: [github.com/intrahealth/client-registry-docs/docs/notebooks/load_bulk_data_in_python.ipynb](https://github.com/intrahealth/client-registry-docs/blob/master/docs/notebooks/simple_query.ipynb)

# In[1]:


#!/usr/bin/env python3
from pathlib import Path
from requests_pkcs12 import get, post
import pandas as pd
import numpy as np

import recordlinkage

import fhirclient.models.patient as p
import fhirclient.models.humanname as hn
import fhirclient.models.contactpoint as cp
import fhirclient.models.fhirdate as fd
import fhirclient.models.identifier as ident
from fhirclient import client

import json
import time
import itertools

# suppress warning: "Certificate for localhost has no `subjectAltName`, falling back to check for a `commonName` for now"
import urllib3
urllib3.disable_warnings(urllib3.exceptions.SubjectAltNameWarning)


# In[2]:


# versions
print("Pandas version: {0}".format(pd.__version__),'\n')
print("Python Record Linkage version: {0}".format(recordlinkage._version.get_versions()['version']),'\n')
print("Numpy version: {0}".format(np.__version__),'\n')
print("FHIR client version: {0}".format(client.__version__),'\n')


# In[3]:


# path to your git clone of github.com/intrahealth/client-registry
crhome = Path.home() / 'src' / 'github.com' / 'intrahealth' / 'client-registry'
clientcert = crhome / 'server' / 'sampleclientcertificates' / 'openmrs.p12'
servercert = crhome / 'server' / 'certificates' / 'server_cert.pem'
csv_file = crhome / 'tests' / 'uganda_data_v21_20201501.csv'


# In[4]:


df_a = pd.read_csv(csv_file)
# df_a = df_a.set_index('rec_id')
print('Number of records :', len(df_a))
print(df_a.head())


# In[5]:


# some cleaning
df_a['rec_id'] = df_a['rec_id'].str.strip()
df_a['sex'] = df_a['sex'].str.strip()
df_a['given_name'] = df_a['given_name'].str.strip()
df_a['surname'] = df_a['surname'].str.strip()
df_a['date_of_birth'] = df_a['date_of_birth'].str.strip()
df_a['phone_number'] = df_a['phone_number'].str.strip()
df_a['uganda_nin'] = df_a['uganda_nin'].str.strip()
df_a['art_number'] = df_a['art_number'].str.strip()

df_a['sex']= df_a['sex'].replace('f', 'female')
df_a['sex']= df_a['sex'].replace('m', 'male')
print(df_a['sex'].value_counts())

# fhirclient validates and some birthdate fields are empty/improperly formatted
# remove non-digits
df_a['date_of_birth'] = df_a['date_of_birth'].str.extract('(\d+)', expand=False)
# force into datetime (coerce has benefit that it removes anything outside of 8 digits)
df_a['date_of_birth'] =  pd.to_datetime(df_a['date_of_birth'], errors='coerce')
# now back into str or fhirdate will complain
df_a['date_of_birth'] = df_a['date_of_birth'].apply(lambda x: x.strftime('%Y-%m-%d')if not pd.isnull(x) else '')

print(df_a.head())


# In[6]:


# default server/path
server = "https://localhost:3000/Patient"
# 3 records, modify if more are required
limit = 100
for index, row in itertools.islice(df_a.iterrows(), limit):
# for index, row in df_a.iterrows():
    patient = p.Patient() # not using rec_id as pandas id, leaving empty
    patient.gender = row['sex']
    name = hn.HumanName()
    name.given = [row['given_name']]
    name.family = row['surname']
    name.use = 'official'
    patient.name = [name]
    phone = cp.ContactPoint()
    phone.system = 'phone'
    phone.value = row['phone_number']
    patient.telecom = [phone]
    patient.birthDate = fd.FHIRDate(row['date_of_birth'])
    emr = ident.Identifier()
    emr.system = 'http://clientregistry.org/openmrs'
    emr.value = row['rec_id']
    art = ident.Identifier()
    art.system = 'http://system1/artnumber'
    art.value = row['art_number']
    nin = ident.Identifier()
    nin.system = 'http://system1/nationalid'
    nin.value = row['uganda_nin']
    patient.identifier = [emr, art, nin]
    # print(json.dumps(patient.as_json()))

    headers = {'Content-Type': 'application/json'}
    start = time.time()
    response = post(server, headers=headers, data=json.dumps(patient.as_json()), 
                    pkcs12_filename=clientcert, pkcs12_password='', verify=servercert)
    end = time.time()
    print(index, response.headers['location'], " | ", round((end - start), 1), "ms") # response.headers['Date']
    # print(response.headers)


# In[ ]:




