# libp2p-kad-dht-trail

Hello world local libp2p network with kad-dht rendezvous discovery


# How to test

First clone the repo, then goto `bootstrap` and run that. You will get a multiaddress.

```
$ cd bootstrap
$ go run main.go

This node is  QmRDzQAZ9bdDWwfK4i33qtPW3Eziditovvg3D7Sn3dWN7f
Available multiaddrs :
/ip4/127.0.0.1/tcp/4000/p2p/QmRDzQAZ9bdDWwfK4i33qtPW3Eziditovvg3D7Sn3dWN7f

```
Now in root run the `main.go` file

```
$ go run main.go --bootstrap /ip4/127.0.0.1/tcp/4000/p2p/QmRDzQAZ9bdDWwfK4i33qtPW3Eziditovvg3D7Sn3dWN7f --port 3000
```

you will be connected to the bootstrap. In another terminal run the same command with different port
```
$ go run main.go --bootstrap /ip4/127.0.0.1/tcp/4000/p2p/QmRDzQAZ9bdDWwfK4i33qtPW3Eziditovvg3D7Sn3dWN7f --port 3001
```

You should see bot nodes getting connected

# Project52

It is one of my [project 52](https://github.com/Sab94/project52).
