try {
  // Initialize config server replica set
  config = {
    "_id": "configrs",
    "members": [
      { "_id": 0, "host": "configsvr1:27019" },
      { "_id": 1, "host": "configsvr2:27020" },
      { "_id": 2, "host": "configsvr3:27021" }
    ]
  };
  rs.initiate(config);
  rs.status();
} catch (error) {
    print("Error occurred during initialization of config server replica set:");
    printjson(error);
}


