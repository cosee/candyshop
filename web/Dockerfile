FROM node

EXPOSE 80

RUN npm install express
RUN npm install redis
RUN npm install http-proxy
COPY static/ /static/
COPY web.js /
CMD ["node", "web.js"]
