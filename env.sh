if [[ $ZSH_EVAL_CONTEXT == "" ]]; then
    echo "Please run with: source env.sh"
else
    export GOPATH="$HOME/Documents/Programming/GoLang/Notes/"
    echo $GOPATH
fi
