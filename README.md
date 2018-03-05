# Login

`docker-compose rm -f && docker-compose up --build`

Example app using go api + postgres db + react app 
https://github.com/McMenemy/GoDoRP


Docker development vs prod
https://docs.docker.com/compose/extends/#extending-services
Docker environment variables
https://docs.docker.com/compose/environment-variables/#the-env-file


Next steps

- move app module to api 
- create app module containing react frontend app
- handle authentication from the react app side
- add app to docker-compose for local development (GoDoRP is a good template to follow)
- add HTTPS support for better security