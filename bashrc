export PROJECT_HOME=$(pwd)
export GOPATH=$PROJECT_HOME/go
export GOBIN=$PROJECT_HOME/go/bin
export PATH=$PROJECT_HOME/bin:$GOBIN:$PROJECT_HOME/web/node_modules/.bin:$PATH
export GO15VENDOREXPERIMENT=1
echo "Rakewire configured"
