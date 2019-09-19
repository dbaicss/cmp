package ci
//package main


import (
	"github.com/bndr/gojenkins"
	"fmt"
	"crypto/tls"
	"net/http"
)

var  jenkins *gojenkins.Jenkins
func BuildExistJob(jobname string) error {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport:tr}

	//初始化jenkins

	jenkins,err := gojenkins.CreateJenkins(client, "https://jenkins", "admin", "xxxxx").Init()

	if err != nil {
		fmt.Printf("Init Jenkins wrong...err:%v\n",err)
		return err
	}
	//job, err := jenkins.GetJob("account-auth")
	job, err := jenkins.GetJob(jobname)
	if err != nil {
		fmt.Printf("Job:%s does not exist\n",jobname)
		return err
	}
	config,err := job.GetConfig()
	if  err != nil {
		fmt.Printf("get congfig failed,err:%v\n",err)
		return err
	}
	fmt.Printf("config:%s\n",config)
	return nil
	/*
	param,err := job.GetParameters()
	//获取parameters,一般都是tag发布的版本号
	res := make(map[string]string)
	for _,p := range param {
		fmt.Printf("params:%#v\n",p.DefaultParameterValue)
		res[p.DefaultParameterValue.Name] = p.DefaultParameterValue.Value.(string)
	}
	//buildId,err := jenkins.BuildJob("account-auth",res)
	//data, err := jenkins.GetBuild("account-auth", buildId)

	buildId,err := jenkins.BuildJob(jobname,res)
	data, err := jenkins.GetBuild(jobname, buildId)
	if err != nil {
		panic(err)
	}

	if "SUCCESS" == data.GetResult() {
		fmt.Println("This build succeeded")
	}
	 */
}

func BuildNewJob(jobName,masterUrl,mvnCmd string)  {
	confString := `config:<?xml version='1.1' encoding='UTF-8'?>
	<maven2-moduleset plugin="maven-plugin@2.17">
	  <actions/>
	  <description></description>
	  <keepDependencies>false</keepDependencies>
	  <properties>
		<com.dabsquared.gitlabjenkins.connection.GitLabConnectionProperty plugin="gitlab-plugin@1.5.11">
		  <gitLabConnection></gitLabConnection>
		</com.dabsquared.gitlabjenkins.connection.GitLabConnectionProperty>
		<jenkins.model.BuildDiscarderProperty>
		  <strategy class="hudson.tasks.LogRotator">
			<daysToKeep>-1</daysToKeep>
			<numToKeep>5</numToKeep>
			<artifactDaysToKeep>-1</artifactDaysToKeep>
			<artifactNumToKeep>-1</artifactNumToKeep>
		  </strategy>
		</jenkins.model.BuildDiscarderProperty>
		<hudson.model.ParametersDefinitionProperty>
		  <parameterDefinitions>
			<net.uaznia.lukanus.hudson.plugins.gitparameter.GitParameterDefinition plugin="git-parameter@0.8.1">
			  <name>develop</name>
			  <description>develop</description>
			  <uuid>e4f500b0-0990-46bb-816e-13b30e98523b</uuid>
			  <type>PT_BRANCH</type>
			  <branch></branch>
			  <tagFilter>*</tagFilter>
			  <branchFilter>.*</branchFilter>
			  <sortMode>NONE</sortMode>
			  <defaultValue></defaultValue>
			  <selectedValue>NONE</selectedValue>
			  <quickFilterEnabled>false</quickFilterEnabled>
			</net.uaznia.lukanus.hudson.plugins.gitparameter.GitParameterDefinition>
		  </parameterDefinitions>
		</hudson.model.ParametersDefinitionProperty>
	  </properties>
	  <scm class="hudson.plugins.git.GitSCM" plugin="git@3.5.1">
		<configVersion>2</configVersion>
		<userRemoteConfigs>
		  <hudson.plugins.git.UserRemoteConfig>
			<url>git@gitlab.icarbonx.cn:icx-demo/icx-blog.git</url>
			<credentialsId>764fe0d9-3511-4584-a80b-610844ffa69f</credentialsId>
		  </hudson.plugins.git.UserRemoteConfig>
		</userRemoteConfigs>
		<branches>
		  <hudson.plugins.git.BranchSpec>
			<name>*/develop</name>
		  </hudson.plugins.git.BranchSpec>
		</branches>
		<doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
		<submoduleCfg class="list"/>
		<extensions/>
	  </scm>
	  <canRoam>true</canRoam>
	  <disabled>false</disabled>
	  <blockBuildWhenDownstreamBuilding>false</blockBuildWhenDownstreamBuilding>
	  <blockBuildWhenUpstreamBuilding>false</blockBuildWhenUpstreamBuilding>
	  <triggers/>
	  <concurrentBuild>false</concurrentBuild>
	  <rootModule>
		<groupId>com.icarbonx.blog.website</groupId>
		<artifactId>blog-web</artifactId>
	  </rootModule>
	  <goals>clean package -U -e</goals>
	  <aggregatorStyleBuild>true</aggregatorStyleBuild>
	  <incrementalBuild>false</incrementalBuild>
	  <ignoreUpstremChanges>false</ignoreUpstremChanges>
	  <ignoreUnsuccessfulUpstreams>false</ignoreUnsuccessfulUpstreams>
	  <archivingDisabled>false</archivingDisabled>
	  <siteArchivingDisabled>false</siteArchivingDisabled>
	  <fingerprintingDisabled>false</fingerprintingDisabled>
	  <resolveDependencies>false</resolveDependencies>
	  <processPlugins>false</processPlugins>
	  <mavenValidationLevel>-1</mavenValidationLevel>
	  <runHeadless>false</runHeadless>
	  <disableTriggerDownstreamProjects>false</disableTriggerDownstreamProjects>
	  <blockTriggerWhenBuilding>true</blockTriggerWhenBuilding>
	  <settings class="jenkins.mvn.FilePathSettingsProvider">
		<path>/usr/share/maven/conf/settings.xml</path>
	  </settings>
	  <globalSettings class="jenkins.mvn.DefaultGlobalSettingsProvider"/>
	  <reporters/>
	  <publishers/>
	  <buildWrappers/>
	  <prebuilders/>
	  <postbuilders>
		<org.jenkinsci.plugins.dockerbuildstep.DockerBuilder plugin="docker-build-step@1.43">
		  <dockerCmd class="org.jenkinsci.plugins.dockerbuildstep.cmd.CreateImageCommand">
			<dockerFolder>$WORKSPACE/</dockerFolder>
			<imageTag>ccr.ccs.tencentyun.com/develop/$JOB_NAME:latest</imageTag>
			<dockerFile>Dockerfile</dockerFile>
			<noCache>false</noCache>
			<rm>false</rm>
			<buildArgs></buildArgs>
		  </dockerCmd>
		</org.jenkinsci.plugins.dockerbuildstep.DockerBuilder>
		<org.jenkinsci.plugins.dockerbuildstep.DockerBuilder plugin="docker-build-step@1.43">
		  <dockerCmd class="org.jenkinsci.plugins.dockerbuildstep.cmd.PushImageCommand">
			<dockerRegistryEndpoint plugin="docker-commons@1.8">
			  <url>https://ccr.ccs.tencentyun.com</url>
			  <credentialsId>dd4fe6e2-7ed1-4ff4-b930-4a32d451ed04</credentialsId>
			</dockerRegistryEndpoint>
			<image>ccr.ccs.tencentyun.com/develop/$JOB_NAME</image>
			<tag>latest</tag>
			<registry></registry>
		  </dockerCmd>
		</org.jenkinsci.plugins.dockerbuildstep.DockerBuilder>
		<org.jenkinsci.plugins.dockerbuildstep.DockerBuilder plugin="docker-build-step@1.43">
		  <dockerCmd class="org.jenkinsci.plugins.dockerbuildstep.cmd.RemoveImageCommand">
			<imageName>ccr.ccs.tencentyun.com/develop/$JOB_NAME:latest</imageName>
			<imageId></imageId>
			<ignoreIfNotFound>false</ignoreIfNotFound>
		  </dockerCmd>
		</org.jenkinsci.plugins.dockerbuildstep.DockerBuilder>
		<hudson.tasks.Shell>
		  <command></command>
		</hudson.tasks.Shell>
	  </postbuilders>
	  <runPostStepsIfResult>
		<name>FAILURE</name>
		<ordinal>2</ordinal>
		<color>RED</color>
		<completeBuild>true</completeBuild>
	  </runPostStepsIfResult>
	</maven2-moduleset>
`
	job,err := jenkins.CreateJob(confString,jobName)
	if err != nil {
		fmt.Printf("create job failed,err:%v\n",err)
		return
	}
	jName := job.Raw.FullName
	fmt.Printf("create job succ,jobName:%s\n",jName)
	err = BuildExistJob(jName)
	if err != nil {
		fmt.Printf("first build job:%s failed,err:%v\n",jobName,err)
		return
	}
	fmt.Printf("first build job:%s succ\n",jobName)
}

//func main()  {
//	BuildExistJob("icx-blog")
//}
