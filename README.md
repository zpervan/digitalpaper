[![Backend Ubuntu](https://github.com/zpervan/digitalpaper/actions/workflows/backend_ubuntu_ci.yml/badge.svg)](https://github.com/zpervan/digitalpaper/actions/workflows/backend_ubuntu_ci.yml)
[![Backend Windows](https://github.com/zpervan/digitalpaper/actions/workflows/backend_windows_ci.yml/badge.svg)](https://github.com/zpervan/digitalpaper/actions/workflows/backend_windows_ci.yml)
[![Frontend Ubuntu](https://github.com/zpervan/digitalpaper/actions/workflows/frontend_ubuntu_ci.yml/badge.svg)](https://github.com/zpervan/digitalpaper/actions/workflows/frontend_ubuntu_ci.yml)
[![Frontend Windows](https://github.com/zpervan/digitalpaper/actions/workflows/frontend_windows_ci.yml/badge.svg)](https://github.com/zpervan/digitalpaper/actions/workflows/frontend_windows_ci.yml)

# Digital Paper #

## Environment ##

### Frontend ###
- Build systems: NPM and YARN
- Typescript

### Backend ###
- Golang 1.19.1

### Third Party Libraries ###
Frontend:
- TODO

Backend:
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Google UUID](https://github.com/google/uuid)

## Setup ##

### Ubuntu ###
1. Follow the [official Docker page](https://docs.docker.com/engine/install/ubuntu/) guide instruction to install it.

## Start the application locally ##
1. Position your terminal in the root of the project
2. Run the command `docker-compose up -d --build` to build and run the containers in detached mode
3. To connect to the web page, type in your browser `localhost:3500`
