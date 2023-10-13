# Transaction Service
### Getting Started
#### Launch the application

    make run    

#### Store a transaction
Using postman or another rest client...

    POST http://localhost:8080/transaction

    {
        "description": "A holiday somewhere nice",
        "transactionDate": "2023-05-01",
        "amountInCents": 100
    }

You will receive a response containing the id of the stored transaction...

    {
        "id": "dfe3adb4-6971-11ee-a606-acde48001122"
    }

#### Fetch a transaction
Specify the id of the transaction to fetch, along with the name of the country (according to the US Treasury Exchange
Rate dataset) of which you would like the transaction amount converted to...

    GET http://localhost:8080/transaction/dfe3adb4-6971-11ee-a606-acde48001122?country=Australia

If you have provided a country for which the dataset has an exchange rate record no older than six months prior to the
transaction date, you will receive a response similar to the following...

    {
        "transaction": {
            "id": "dfe3adb4-6971-11ee-a606-acde48001122",
            "description": "A holiday somewhere nice",
            "transactionDate": "2023-05-01",
            "amount": {
                "usdAmountInCents": 100,
                "convertedAmountInCents": 154,
                "exchangeRate": 1.542
            }
        }
    }

### Context Diagram


                       │
                       │ json/http
                       │
          ┌────────────▼────────────┐
          │                         │
          │   transaction-service   │
          │                         │
          └────────────┬────────────┘
                       │
                       │
                       │ json/http
                       │
                       │
          ┌────────────▼────────────┐
          │                         │
          │      us treasury        │
          │    exchange rate api    │
          │                         │
          └─────────────────────────┘






### Assumptions

1. The client of the transaction service expects to interact with it via an api.
2. The client will call the api in a typical synchronous manner.  (as opposed to an async event based model) 
3. The client is typical in that it interacts with apis using JSON.
4. Transactions received with the same details are not identical.  It is feasible that multiple transactions with the same description, date and amount are received for different transaction events.  In a production system, consider including a unique identifier in the request so that the back end can store the transaction in an idempotent manner without risk of duplication.  This is important in a distributed system.  Alternatively, consider including fields in the request from which a natural key can be formed.
5. The maximum transaction amount the system needs to support, including in its calculations is well within the bounds of safe integer values.
6. The transaction date received will be in UTC timezone.
7. The transaction date must be today or in the past.  It doesn't seem to make sense to have the system handle future purchases, but that would be something to confirm. 
8. A transaction description must be at least one character long.
9. A transaction amount cannot be zero, since this doesn't seem to make sense as a transaction.
10. Investigated if countries can have multiple currencies.  China apparently has, but the Exchange Rate API returns only one currency for China, so it is safe to assume all countries have one currency - for our purposes.
11. Investigated if any countries use the US dollar as their official currency.  There are a few (e.g. Micronesia), but the Exchange Rate API returns an exchange rate of 1.0 for those.  So it is safe to assume this will not cause any complications.
12. Transactions are allowed to have a negative amount, potentially to support refunds etc.  Requirements such as this would need to be confirmed.
13. A country name will not be shorter than two characters long.
14. A country name will not be longer than 255 characters.
15. During normal operation, it will not take longer than 5 seconds to get a response from the Treasury API
 
### Scope

Due to time constraints, the following have been considered out of scope:

* Authentication / Authorisation
* Input Sanitization
* Caching
* Packaging as a container for distribution
* Graceful shutdown of the http server
* Externalised and environment specific configuration
* Correlation IDs for distributed tracing
* Context Cancellation


### Notes
* All user input would ideally be sanitised using something like [bluemonday](github.com/microcosm-cc/bluemonday), although I haven't implemented this due to time constraints.
* It would be a good idea to cache responses from the Treasury API so that we are not hammering it under volume.  This could be implemented at the `forex.Repository` layer.
* I have made sure to set `MaxConnsPerHost` in the http client so that connection pooling settings are not restrictive.  This would need to be tuned properly in production.
* Using go standard library logger.  In a production system, consider using a more fully functional logger such as [Zerolog](https://github.com/rs/zerolog), [Zap](https://github.com/uber-go/zap), or [Apex](https://github.com/apex/log). 
* We are using an in memory repository to store transactions.  In a production system this simple approach would not likely be viable as it does not provide long term storage.
* Using simple R/W mutex to perform synchronisation on the in memory map used for storage.  In a production system this approach may / may not be performant, although in that scenario a real database would likely be used.
* Validation frameworks can help achieve code consistency, however, they can also introduce constraints.  Since I have relatively lean experience with the gin validation library I opted to stick with a simple custom implementation so that I was not subjected to any such constraints.   In this instance, using the gin validation framework would be the most obvious option, however it would be worth evaluating various other validation options before committing.
* For integration testing I have opted not to tightly integrate my testing with the gin framework.  That approach is a valid option which would reduce the setup code, however, it also couples more things to gin.
* The port used by the web server is configured to be dynamic when run via integration tests.  This means we can have the server running locally on port `8080` and it won't interfere with the running of integration tests.  Also on a CI build agent there should never be port conflicts related to the integration tests.
* I have hand-crafted mocks, but in general would usually use [mockery](https://github.com/vektra/mockery) to generate them.
* It might be sensible to add business validation to help make sure amounts received / returned do not go out of bounds of the numeric data types used.  Although I haven't implemented that.
* I haven't used it here, but consider use of [lightweight architecture decision records](https://github.com/peter-evans/lightweight-architecture-decision-records) to help retain context and provide it for self reference and that of engineers new to the project.
* UUIDs have been used for generated IDs as they are effectively unique and do not require synchronisation to generate. e.g. sequential ids require that we know what the previous id was.
* A sequential ID generator has been used for testing purposes and is only wired in for tests.  Outside of testing, the UUID generator is used.


### Development
#### Run all tests
    make test
#### Run unit tests
    make test unit
#### Run integration tests
    make test-integration
#### Launch from code
    make run
#### Build the executable
    make build
