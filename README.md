# Phone and Email verificator

This project intend to access the [API](https://www.bigdatacloud.com/phone-email-verification/phone-number-validation-api#endpoint) to validate phones from [BigData](https://www.bigdatacloud.com/) website, to collect the responses and storage in a MongoDb.


## Prerequisites

To run this project, you must have intalled:

* You must have a KEY from BigData website, do a SignUp om the website and create one.
* docker-compose
* docker
* go (1.20 or higher)

## To run the project
1) Access the project folder through a shell.

2) Export your API_BDC_KEY with the command:

    `export API_BDC_KEY=<YOUR_KEY>`

3) In one shell, connect with the MongoDb running the following docker-compose command on the project folder root:

    `docker-compose up`
    * Sometimes the DB takes a realtive long time to get up.

4) In another shell, run the app:

    `go run cmd/phoneEmailVerification/phoneEmailVerification.go`

5) If everything runs ok, you should see this message on the shell:

    `> api server running on http://localhost:8080`

## Available routes
### Verify phones:

* (POST) Search and save a number:
    * Example:

        `http://localhost:8080/api/phoneNumber/2018675309/countryCode/us/localityLanguage/en`

* (GET) Get specific number already collected on DB:
    * Example:

        `http://localhost:8080/api/getPhone/2018675309`

* (GET) Get all numbers already collected:
    * Example:
    
        `http://localhost:8080/api/getAllPhones`

### Verify emails:
* Not available =/