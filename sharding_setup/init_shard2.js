try {
  // Initialize shard2 replica set
  config = {
    "_id": "shard2rs",
    "members": [
      { "_id": 0, "host": "shard2svr1:27025" },
      { "_id": 1, "host": "shard2svr2:27026" }
    ]
  };
  rs.initiate(config);
  rs.status();

  // Wait for the shard2 replica set to initiate
  sleep(5000);
} catch (error) {
    print("Error occurred during initialization of shard2 replica set:");
    printjson(error);

}

