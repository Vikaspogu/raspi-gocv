node {
    checkout scm

    docker.withRegistry('docker.io', 'docker-auth') {
        def customImage = docker.build("vikaspogu/rpi-node-cm")
        customImage.push()
        customImage.push('latest')
    }
}
