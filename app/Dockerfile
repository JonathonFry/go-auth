FROM node

# if left blank app will run with dev settings
# to build production image run:
# $ docker build ./app --build-args app_env=production
ENV NPM_CONFIG_LOGLEVEL warn

RUN mkdir -p /app
WORKDIR /app
COPY ./ ./

RUN npm install

# if dev settings will use create-react start script for hot code relaoding via docker-compose shared volume
# if production setting will build optimized static files and serve using http-server
CMD if [ ${ENV} = production ]; \
	then \
	npm install -g http-server && \
	npm run build && \
	cd build && \
	hs -p 3000; \
	else \
	npm run start; \
	fi