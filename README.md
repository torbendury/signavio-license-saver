# signavio-license-saver

This is a simple tool for license rotation in Signavio. It retrieves all users from Signavio and deletes them. This frees up license seats.

*It might be against the terms of service of Signavio to actually use this tool, so I would advise against it. It is meant as a proof of concept and for educational purposes only.*

- [signavio-license-saver](#signavio-license-saver)
  - [Use Case](#use-case)
  - [Usage](#usage)
  - [Development](#development)
  - [Signavio APIs](#signavio-apis)
  - [License](#license)
  - [Project](#project)

## Use Case

The use case for this tool is to free up license seats in Signavio. If you have a lot of users and you want to free up some seats, you can use this tool to delete all users and then let Signavio re-assign licenses to the users actually logging in and using the tool.

## Usage

The tool does not need any external dependencies, so for compiling you only need `make` and `go` v1.21.5 installed.

Clone the repository and install the binary yourself:

```bash
make install

signavio-license-saver -h
```

`signavio-license-saver` offers a few flags to customize behavior.

Optionally, if you're into `dotenv` files, you can create a `.env` file with the following content:

```bash
TENANT=****
URL=****
USER=****
PASSWORD=****
ALLOWLIST=important@mail.com,veryimportant@mail.com
```

And afterwards run `make run` to execute the tool. Make will automatically load the `.env` file and use the environment variables from there.

## Development

This repo is built to be extendable, you can enhance the API at any time. Under `/hack` you can find a mock server which you can use to test changes. It implements some basic routes.

## Signavio APIs

Apart from the documented APIs, there is one undocumented API which is used by this tool. It is the `GET /p/user` endpoint. It returns a list of all users in the Signavio tenant. I retrieved it from my browser journey. **It is not guaranteed to work in the future.**

The other APIs are documented and can be found in the official SAP Signavio docs.

## License

This project uses GPLv3.0. You can find the license [here](LICENSE).

## Project

Feel free to contribute to this project! You can open issues, create pull requests or just fork the project and use it for your own purposes. I am happy to receive any feedback.

I'm not affiliated with Signavio or SAP in any way. This is a personal project and not an official tool.

If you find bad practices, bugs or anything else, please let me know. I'm happy to learn and improve.

Made with ❤️ by [Torben](https://torbendury.de)
