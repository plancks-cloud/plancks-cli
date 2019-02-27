package model

//AppPort is the port that the brutus application will listen on
const AppPort = 6227

//ListenPort is the port that the app listens to internally
const ListentPort = 6227

//InstallMaxHealthChecks is the number of times it will check for the health service during install
const InstallMaxHealthChecks = 5

//InstallSleepBetweenChecksis is the number of seconds the app will wait between checking for the health service during install
const InstallSleepBetweenChecks = 5
