const request = require('request');

request(JSON.parse(process.argv[2]), (err, response, body) => {
  if (err) {
    console.error(err);
    return;
  }
  console.log(`Status Code: ${response.statusCode}`);
  console.log(`Body: ${body}`);
});
