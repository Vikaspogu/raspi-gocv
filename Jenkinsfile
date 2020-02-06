pipeline {
  agent {
    kubernetes {
      yamlFile 'kubernetes/jenkinsPod.yaml'
    }
  }
  stages {
    stage('build image') {
      step{
        checkout scm
        container('docker') {
            sh 'DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1  docker run --rm --privileged multiarch/qemu-user-static --reset -p yes'
            sh 'cd `pwd` && DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1 docker build --platform linux/arm64 -t "docker.io/vikaspogu/rpi-node-cm" .'
        }
      }
    }
    stage('push image') {
      step{
        container('docker') {
          sh 'DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1 docker push docker.io/vikaspogu/rpi-node-cm'
        }
      }
    }
    stage('deployment'){
      step{
        container('kubectl'){
          sh 'kubectl get pods'
        }
      }
    }
  }
}
