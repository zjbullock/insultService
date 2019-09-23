# insultService
A service to send funny, yet charming insults.

This project was designed for both rusty bot and my personal (and eventual) insult app.

This project utilizes Go, FireStore, and Graphql.

Graphql provides scalability in that data can be controlled by limiting responses to exactly what is necessary, no more and no less.

FireStore is utilized to store all possible words (adjectives, nouns, verbs) that can be chosen, and combinations that are actually used.

![alt text](https://raw.githubusercontent.com/zjbullock/insultService/master/insult_demo.PNG)
