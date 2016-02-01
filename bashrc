export PROJECT_HOME=$(pwd)
export GOPATH=$PROJECT_HOME/go/vendor:$PROJECT_HOME/go
export GOBIN=$PROJECT_HOME/go/vendor/bin
export PATH=$PROJECT_HOME/bin:$GOBIN:$PROJECT_HOME/web/node_modules/.bin:$PATH
echo "Rakewire configured"
