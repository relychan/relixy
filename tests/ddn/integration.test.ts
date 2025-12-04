// these tests expect the Docker Compose services in `/tests/docker-compose.yaml` to be running

test("We can make a request using query parameters", async () => {
  const url = "http://localhost:3280/v1/api/rest/artistbyname/Queen";
  await fetch(url).then(async (response) => {
    const json = await response.json();
    expect(json).toMatchSnapshot();
  });
});

test("We can make a request using query parameters", async () => {
  const url = "http://localhost:3280/v1/api/rest/artists?limit=10&offset=20";
  await fetch(url).then(async (response) => {
    const json = await response.json();
    expect(json).toMatchSnapshot();
  });
});
