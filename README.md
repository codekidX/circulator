# Circulator - An embedded config management server


This project is just a random thought-experiment and a 10mins-project and sort of like a proof of concept that it is doable.


The how of it is really interesting. I was just browsing through some tutorials about how to setup configs in python and how to maintain them, it just struck me that I already have applications running in Node and Go which has it's own config management thingy implemented and I just thought do I always need to perform this symphony whenever I write a program in a different language?

The answer is/was a big fat **NO** !!


## What?

Well it is written in Go for a reason and mainly because it's cross platform nature and performance with respect to the HTTP server. 

You can implement the server in your own way but this is how I serve configs to my multiplatform/multilingual applications now. 

## How circulator works?

It embeds the static configuration JSON5 inside the final binary that it builds and it serves a response as JSON to the applications that requires them.

The `config/` folder is the heart of the circulator and it holds all the JSON5 configuration. The heirarchy that works for me is application based, so each JSON5 file belongs to a single application.

## Advantages of this approach

- You don't have any configuration files in your final deploy
- You don't need to use separate library for application _(of any language)_ like .env, and/or external dependencies that wants the configs
- You have a single secret string that could be passed through environment variables
- The target application that has the configs have it in a volatile state

## Accessing configs

- API call to `localhost:$YOUR_PORT/?app=$CONFIG_FILE_NAME` should give you a JSON response of the file.
- Through **dot notation** - `localhost:$YOUR_PORT/?app=$CONFIG_FILE_NAME.$KEY` will give you the value of nested object
> _Example: `coolapp.auth.token` will return "mytoken"_ if config is
> ```
> {"auth": {"token": "mytoken"}}
> ```

**$YOUR_PORT** can be defined inside `./config/__cconfig.json`

## What is __cconfig.json ?

It is the configuration file for the behaviour of circulator.

| Key | Type | Description |
|------|------|-----------|
| port | number | defines on which port circulator should run |
| secret | string | the bearer token that you can set for authorizing requests |
| protect | boolean | protect server from being used by other people and ask for secret when your server runs |

## Running

```sh
go run main.go
```

## Error codes

- `400`: not enough params, missing `app` maybe?
- `401`: bearer token wrong or missing
- `404`: config file not found
- `500`: wrong dot notation, unable to parse JSON

## TODOs

- [x] Add password protection to the binary so that no one else in the world can run it
- [ ] Write tests
- [x] Support for accessing through dot notation
- [ ] Spit out config files from the binary
- [ ] Support for different environment - this is a long shot and maybe it should not be included to make this a complex application _(because we could include "$ENV" key and same object to a single file the problem is it'll become big and messy)_


## Confessions

- It is not a great tool if you are working as a team - Have any idea that can make this project collaborative - [ping me](https://github.com/codekidX/circulator/issues). 
> But it becomes one if you are using it in a private environment, just remove the `config/` from `.gitignore`
- It does not support other configuration languages.
- Why JSON5? - Comments, because I have a 3 months threshold of remembering code and also it excels at single way of writing things which is a rare entity these days.


## License

This application is licensed under `Buy Author A Pizza` 🍕 license.