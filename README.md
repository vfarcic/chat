```bash
sudo docker run -d --name oauth \
	-p 8080:8080 \
	-v /etc/ssl/certs:/etc/ssl/certs \
	-v /data/oauth:/data/db \
	-e SEC_KEY="Bla" \
	-e REDIRECT_URL="http://localhost:8081/auth/google/callback" \
	-e GOOGLE_CLIENT_ID="472858977716-ej3ca5dtmq4krl7m085rpfno3cjp2ogp.apps.googleusercontent.com" \
	-e GOOGLE_SECRET="OnkptU4BTdE45mi-b3hACdAY" \
	-e GOOGLE_REDIRECT_URL="http://localhost:8080/auth/google/callback" \
	vfarcic/oauth

bower install

go build -o chat && ./chat -addr ":8081"
```

Make sure that authName (mandatory) and authAvatarURL (optional) are created before user enter a room.