Explanation
===========

* Make sure that authName (mandatory) and authAvatarURL (optional) cookies are created before user enters a room.
Optionally, use vfarcic/oauth to authenticate and create those cookies.
* Submit Chat and Display Chat Web components can be imported and displayed on a page.

Compile
=======

```bash
bower install

sudo docker run --rm \
	-v $PWD:/usr/src/chat \
	-v $GOPATH:/go \
	-w /usr/src/chat \
	golang:1.4 \
	go get -d -v && go build -v -o chat go/*.go
```

Build Docker Container
======================

```bash
sudo docker build -t vfarcic/chat .
```

Run
---

```bash
sudo docker run -d --name chat \
	-p 8080:8080 \
	-v /etc/ssl/certs:/etc/ssl/certs \
	-e PORT=8080 \
	vfarcic/chat
```

Embed "Chat" Web Components
===========================

```html
<html>
<head>
	<!--Import Required Polymer Components-->
    <link rel="import" href="../bower_components/paper-styles/classes/global.html">
    <link rel="import" href="../bower_components/polymer/polymer.html">
    <link rel="import" href="../bower_components/iron-icon/iron-icon.html">
    <link rel="import" href="../bower_components/iron-icons/iron-icons.html">
    <link rel="import" href="../bower_components/paper-button/paper-button.html">
    <link rel="import" href="../bower_components/paper-fab/paper-fab.html">
    <link rel="import" href="../bower_components/paper-input/paper-textarea.html">
    <link rel="import" href="../bower_components/paper-toast/paper-toast.html">
    <link rel="import" href="../bower_components/paper-material/paper-material.html">
    <link rel="import" href="../components/chat/submit-chat.html">
    <link rel="import" href="../components/chat/display-chat.html">
    <!--Import "Chat" Components-->
	<link rel="import" href="http://localhost:8080/components/chat/submit-chat.html">
	<link rel="import" href="http://localhost:8080/components/chat/display-chat.html">
</head>
<body>
	<!--Display "Chat" Components-->
	<display-chat></display-chat>
	<submit-chat></submit-chat>
</body>
```