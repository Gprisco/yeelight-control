# Yeelight Control
A simple golang cli program that implements the [Yeelight Inter-Operation Specification](https://www.yeelight.com/download/Yeelight_Inter-Operation_Spec.pdf) for executing actions on lightbulbs under the same local network.

## Operations
### Search
Our program will multicast a search request and wait for a unicast response to be received.

## To do
* Handle multiple unicast **search** responses to be received, every yeelight device available will send one with its location details
* Implement other yeelight operations
