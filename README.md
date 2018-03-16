# go-auth

Sample website to learn modern web technologies.

Built using docker-compose for easier development/deployment. **app|api|db** each have their own Dockerfiles and are able to be run independently.

- api is developed using go + [gorilla mux](https://github.com/gorilla/mux)
- db is a [postgres](https://www.postgresql.org/) SQL db
- app is a [react](https://reactjs.org/) + [redux](https://redux.js.org/) web app

To start run:
`docker-compose up --build`

To run with a clean db use:
`docker-compose rm -f && docker-compose up --build`
This removes the existing db container and starts fresh

## Next steps
- Better validation for JWT token. Ensure that user is still valid
- Autorefresh JWT token so user doesn't have to re-login
- add HTTPS support for better security using [Let's Encrypt](https://letsencrypt.org/)


## Shoutouts
React tutorials:
- https://thinkster.io/tutorials/react-redux-ajax-middleware/
- http://jasonwatmore.com/post/2017/09/16/react-redux-user-registration-and-login-tutorial-example
- https://auth0.com/blog/secure-your-react-and-redux-app-with-jwt-authentication/

Go architecture:
- [Example app using go api + postgres db + react app](https://github.com/McMenemy/GoDoRP)
- [Docker development vs prod](https://docs.docker.com/compose/extends/#extending-services)
- [Docker environment variables](https://docs.docker.com/compose/environment-variables/#the-env-file)
