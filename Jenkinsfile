pipeline {
  agent {
    kubernetes {
      yamlFile 'kubernetes/jenkinsPod.yaml'
    }
  }
  stages {
    stage('build image') {
      steps{
        checkout scm
        container('docker') {
            sh 'DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1  docker run --rm --privileged multiarch/qemu-user-static --reset -p yes'
            sh 'cd `pwd` && DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1 docker build --platform linux/arm64 -t "docker.io/vikaspogu/rpi-node-cm" .'
        }
      }
    }
    stage('push image') {
      steps{
        container('docker') {
          sh 'DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1 docker push docker.io/vikaspogu/rpi-node-cm'
        }
      }
    }
    stage('deployment'){
      steps{
        container('kubectl'){
          sh 'kubectl get pods'
        }
      }
    }
  }
}
