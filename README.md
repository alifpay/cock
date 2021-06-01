# Introduction to Distributed SQL and CockroachDB

- A distributed SQL database is a single relational database which replicates data across multiple servers. 
 - Google's Spanner popularized the modern distributed SQL database concept. Google described the database 
 and its architecture in a 2012 whitepaper called "Spanner: Google's Globally-Distributed Database." 
 - Distributed SQL databases have the following general characteristics:
    - synchronous replication
    - strong transactional consistency across at least availability zones (i.e. ACID compliance) 
    - relational database front end structure – meaning data represented as tables with rows and columns similar to any other RDBMS
    - automatically sharded data storage
    - underlying key–value storage
    - native SQL implementation
    
