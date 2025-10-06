### Build UI
```
cd webui
npm run build
tar czvf dist.tar.gz dist/
```

### Deploy on nginx
```
apt install nginx

cat > /etc/nginx/conf.d/myserver.conf << EOF
server {
    listen 192.168.100.100:9090;
    server_name myhealth.com;

    location / {
        root /var/www/html/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
}
EOF

mv dist.tar.gz /var/www/html
cd /var/www/html
tar xzvf dist.tar.gz

systemctl restart nginx
```