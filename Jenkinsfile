podTemplate(yaml: '''
apiVersion: v1
kind: Pod
spec:
  nodeSelector:
    k3s.io/hostname: hp-mini
  containers:
  - name: docker
    image: docker:19.03.1
    command:
    - sleep
    args:
    - 99d
    volumeMounts:
      - name: jenkins-docker-cfg
        mountPath: /root/.docker
    env:
      - name: DOCKER_HOST
        value: tcp://localhost:2375
  - name: docker-daemon
    image: docker:19.03.1-dind
    securityContext:
      privileged: true
    env:
      - name: DOCKER_TLS_CERTDIR
        value: ""
  volumes:
  - name: jenkins-docker-cfg
    projected:
      sources:
      - secret:
          name: regcred
          items:
            - key: .dockerconfigjson
              path: config.json
''') {
  node(POD_LABEL) {
    stage('build image') {
      container('docker') {
          checkout scm
          sh 'export DOCKER_BUILDKIT=1 && docker build --platform=local -o . git://github.com/docker/buildx && mv buildx ~/.docker/cli-plugins/docker-buildx'
          sh 'cd `pwd` && docker buildx build --platform linux/arm64 -t "docker.io/vikaspogu/rpi-node-cm" .'
      }
    }
    stage('push image') {
      container('docker') {
          sh 'DOCKER_BUILDKIT=1 docker push docker.io/vikaspogu/rpi-node-cm'
      }
    }
  }
}
