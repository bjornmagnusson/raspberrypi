node(pi1) {
    docker.image('resin/rpi-raspbian').inside {
        checkout scm

        stage 'Commit'
        sh 'go get -v github.com/stianeikeland/go-rpio'
        sh 'go get -v github.com/kidoman/embd'
        sh 'go build -v'

        stage 'Test'
        sh 'go test -v'
    }
}
