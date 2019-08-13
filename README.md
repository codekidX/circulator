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

A simple API call to `localhost:$YOUR_PORT/?app="coolapp"` should give you a JSON response.

> $YOUR_PORT can be defined inside `./configs/__cconfig.json`


#### __cconfig.json

Has 2 values: 

- port = defines on which port circulator should run
- secret = the bearer token that you can set for authorizing requests

```json
{
    "port": 8000,
    "secret": "blahblah"
}
```

## Running

```sh
go run main.go
```

## TODOs

- [ ] unchecked
- [ ] Write tests
- [ ] Support for accessing through dot notation
- [ ] Spit out config files from the binary


## Confessions

- It is not a great tool if you are working as a team - Have any idea that can make this project collaborative - [ping me](https://github.com/codekidX/circulator/issues).
- It does not support other configuration languages.
- Why JSON5? - Comments, because I have a 3 months threshold of remembering code and also it excels at single way of writing things which is a rare entity these days.


## License

This application is licensed under `Buy Author A Pizza` 🍕 license.