
## Run the application

To run the next API please follow these steps:

1- Clone the repository and then rename the .env(EXAMPLE) file with the corresponding credentials.
2- You need to create a database called maze, and two collection called spots and quadrants.
3- You can create the quadrants via API but aware that just four quadrants must exist for this reason the attribute type represents the quadrant type that is a enum: TOP_RIGHT, TOP_LEFT, BOTTOM_RIGHT, BOTTOM_LEFT.


Request example:
```json
{
    "type": "TOP_RIGHT",
    "start_point": {
        "x": 0,
        "y": 0
    },
    "limit_point": {
        "x": 25,
        "y": 25
    }
}
```
4- Once you have created them and put credentials, you can make the following command to start the server:

```bash
    $ go run main.go
```

5- You can run check each route in postman, here you are the collection requests:

https://www.getpostman.com/collections/5ff52a517f09b3f36efe

## Tests
To run unit test you can run the follow commands to do it:

```bash
    $ make test name=TestSpot_CreateListDelete
```

Note: just change the test name if you want to test another.

## Golangci Lint
To check run the lint and check what errors we have please run the following command:

```bash
    $ make lint pkg=routes
```

Note: If you want to check another package just change routes, e.g `make lint pkg=repository`