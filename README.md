# dynamic-consent

# 1. outline

The dynamic-consent repository is a research result of blockchain and channel-based Personal Health Record(PHR) platform development supported by KHIDI(Korea Health Industry Development Institute.)
It contains authorization, dynamic consent, medical information sharing records features with Hyperledger Fabric version 2.2 for PHR management using blockchain technology.
Chaincode with GO language and the server using Node.js on this repository are not all components to build the entire PHR platform. Only partial blockchain functions on the PHR platform are open, and also the platform centralization system and the hospital link system are not involved in the repository.

* The repository was not made for the purpose of Hyperledger Fabric education but with an assumption that the readers are senior developers  who know how to analyze and use the framework with basic knowledge of Hyperledger Fabric version 2.2 
* Hyperledger Fabric has various operation processes by each version, so you should know that using other versions may affect regular operation.


# 2. source code description

Path: dynamic-consent/hyperledger_fabric/application
Source codes using Javascript and Node.js run functions developed by blockchain chain-code with Restful API.
- .js files are the first version of source code implemented from basic designs (one upload in 2020)
- v1.0.js files are the final version of the source code implementing all functions (one upload in 2021)

Path: dynamic-consent/hyperledger_fabric/chaincode
Source codes with GO language run functions like authorization, dynamic consent, and medical information sharing records.
- .js files are the first version of source code implemented from basic designs (one upload in 2020)
- v1.0.js files are the final version of the source code implementing all functions (one upload in 2021)
