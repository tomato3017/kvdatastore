# kvdatastore
#### Purpose

KVDatastore or Key-Value Datastore is a data storage abstraction that is used to abstract the data source used to back the 
kv-like map. Its meant to simplify the use of a simple data structure that is backed with a variety of datasources.

#### Supported Datastores
* Memory only(basically a map, useful for testing)
* Afero-backed filessystem
  * Memory
  * OS
  * Union
  * SFTP
  * Anything Afero supports
* Redis
