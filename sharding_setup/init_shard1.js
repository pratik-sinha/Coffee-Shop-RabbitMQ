

try {
  // Initialize shard1 replica set
  config = {
    "_id": "shard1rs",
    "members": [
      { "_id": 0, "host": "shard1svr1:27022" },
      { "_id": 1, "host": "shard1svr2:27023" }
    ]
  };
  rs.initiate(config);
  rs.status();

  // Wait for the shard1 replica set to initiate
  sleep(5000);
} catch (error) {
    print("Error occurred during initialization of shard1 replica set:");
    printjson(error);
  }


