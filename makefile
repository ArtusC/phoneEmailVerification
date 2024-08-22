race ?=
testfile ?=
testname ?=
run_compose := docker-compose

git:
	git add .
	git commit -m "$m"
	git push

# make createDBTest
createDBTest:
	@echo "Creating the mongodb to running the integration tests"
	$(run_compose) -f docker-compose-test.yaml up 

# make createDBApp
createDBApp:
	@echo "Creating the mongodb to running the application"
	$(run_compose) up 

# run unitary tests
# to run all tests, just run: make unitCheck
# to run a specific test file, example: make unitCheck testfile=./context/v3/
# to run a specific test name, example: make unitCheck testfile=./context/v3/ testname=TestServer
# if you want to run these tests with a race flag, put -race=race, example: make unitCheck testfile=./context/v3/ testname=TestServer -race=race
unitCheck:
	./hack/builder-check.sh testfile=$(testfile) testname=$(testname) race=$(race)

# run integration tests
# to run all tests, just run: make integrationCheck
# to run a specific test file, example: make integrationCheck testfile=./context/v3/
# to run a specific test name, example: make integrationCheck testfile=./context/v3/ testname=TestServer
integrationCheck:
	./hack/builder-check-integration.sh testfile=$(testfile) testname=$(testname)