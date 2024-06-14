package main

import (
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func process(source Source) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	for {
		<-ticker.C

		if source.Remote {
			remoteFile(
				source.KeyPath,
				source.User,
				source.Host,
				source.Port,
				source.Name,
				source.Paths,
			)
		} else {
			localFile(source.Name, source.Paths)
		}
	}
}

func localFile(name string, paths []Path) {
	if err := os.MkdirAll(storageDir+"/"+name, os.ModePerm); err != nil {
		log.Print(err)
		return
	}

	for _, path := range paths {
		bytes, err := os.ReadFile(path.BasePath + "/" + path.Filename)
		if err != nil {
			log.Print(err)
			continue
		}

		file, err := os.OpenFile(
			storageDir+"/"+name+"/"+path.Filename,
			os.O_TRUNC|os.O_WRONLY|os.O_CREATE,
			0644,
		)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()

		if _, err := file.Write(bytes); err != nil {
			log.Print(err)
			continue
		}
	}
}

func remoteFile(
	keyPath string,
	user string,
	host string,
	port string,
	name string,
	paths []Path,
) {
	pemBytes, err := os.ReadFile(keyPath)
	if err != nil {
		log.Print(err)
		return
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Print(err)
		return
	}

	sshConfig := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":"+port, &sshConfig)
	if err != nil {
		log.Print(err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Print(err)
		return
	}
	defer session.Close()

	reader, err := session.StdoutPipe()
	if err != nil {
		log.Print(err)
		return
	}

	if err := os.MkdirAll(storageDir+"/"+name, os.ModePerm); err != nil {
		log.Print(err)
		return
	}

	for _, path := range paths {
		file, err := os.OpenFile(
			storageDir+"/"+name+"/"+path.Filename,
			os.O_TRUNC|os.O_WRONLY|os.O_CREATE,
			0644,
		)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()

		if err := session.Start("cat " + path.BasePath + "/" + path.Filename); err != nil {
			log.Print(err)
			continue
		}

		if _, err := io.Copy(file, reader); err != nil {
			log.Print(err)
			continue
		}
	}

	if err := session.Wait(); err != nil {
		log.Print(err)
		return
	}
}
