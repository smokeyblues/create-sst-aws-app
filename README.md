# Serverless Full-Stack Application Generator

This command line tool generates a fullstack serverless application written in typescript and hosted on AWS using the SST framework. The template is based on a simple note taking app called Scratch. The steps for creating this app from scratch (pun not intended) is covered in the SST Guide. I recommend going through the guide at least once in order to become familiar with the framework and the intracacies of this app.

Before using this tool, make sure that you have:

    * an AWS account set up
    * created a IAM user
    * configured the AWS CLI

To run this locally make sure you have Go installed on your system. I know that might be obvious but I tried to show this to a friend on my laptop and these are all the steps I had to go through while trying to do that demonstration before I could actually get it to work. 

Now that we have all that out of the way, let's get into how to use this tool to create your next fullstack serverless app.

    1. click **Use this template > Create new repository**
    2. Give your repository a name and hit the **Create Repoository** button
    3. run the following commands:


```
git clone <REPO_URL>
cd <REPO_NAME>
go run .
```

From here you should encounter the text-based user interface:

![Text-based user interface](/assets/tui-index-screenshot.png)

There are currently two templates to choose from. The usage-based pricing template has a stripe-integration built in. To use this one, make sure that you have saved your stripe secret key using the command:

```
npx sst secret set StripeSecretKey <YOUR_STRIPE_SECRET_TEST_KEY>
```

[Instructions for setting up a stripe account and retrieving your secret key.](https://guide.sst.dev/chapters/setup-a-stripe-account.html)

The no-fee version does not use any stripe integration.

Select your option and hit enter.

You will then be asked to enter the name of your project:

![Enter project name](/assets/enter-demo-project-name.png)

Here I have chosen the name 'demo-project'. After you type in your project name, hit enter and your project will be created in a directory with the same name you chose for your project. From here run the following commands:

```
cd <YOUR-PROJECT-NAME>
npm install
```

Once you have everything installed you can deploy a development server for your project with the following command:

```
npx sst dev
```

## If everything worked as expected, you should see a screen that looks something like this:

![Successfully deployed app in dev stage](/assets/successful-dev-stage-deployment-screenshot.png)

As you can see, several AWS resources have been deployed. These include an API gateway, an identity pool, and a user pool. It has also deployed an S3 bucket for hosting the frontend code serverlessly as well as a S3 bucket that allows users of your app to save files associated with each of their "notes" in this demo app. 

## Hit the down button to the frontend deployment and find the url of your app:

![frontend deployment screen](/assets/frontend-deployment-screen.png)

## And open that location in your browser to find your app running: 
![dev app in browser](/assets/app-in-browser.png)

You can now edit the code in packages/frontend/src and see your changes immediately take effect in the browser. For example, head to /packages/frontend/src/containers/Home.tsx and change the name of the app in the RenderLander() function to the name of your project.

## Your project name should now be front and center:
![Title changed to demo-project](/assets/title-changed.png)

This command-line tool simplifies the process of creating new full-stack serverless applications built with React, TypeScript, Vite on the frontend, and SST and AWS on the backend.  Choose from pre-built templates to kickstart your next project.

These templates have only followed the SST guide up to the point just before the chapters before it goes into the details about [how to get production ready](https://guide.sst.dev/chapters/getting-production-ready.html) and how to use custom domains for your project so be sure to review those chapters once you are happy with the changes you have made to your app. 

If you have no intentions of using this app after you deploy it make sure that you remove all the resources with the following command:

```
npx sst remove
```

This command should destroy all AWS resources created for this project so you don't have a bunch of orphaned elements cluttering up your AWS environment. 

## A note on resources used in this project:

This project is my submission for the final project in CS50x and my first project in Go. Since Go is a google project, I decided to use a google llm as my coding assistant. The model I used was the LearnLM 1.5 Pro Experimental model that can be found in the model options in the [aistudio.google.com/prompts](https://aistudio.google.com/prompts) console. Just hit Create new prompt and you should be able to find the model options on the right. 

I used the model mostly as a search tool to help me quickly familiarize myself with the Go language and some of the libraries I was interested in exploring in this space.

SST uses Go to write their CLI tools so they were the inspiration for this project. I have always been impresses with their CLI tools and found that they were using a library called [bubbletea](https://github.com/charmbracelet/bubbletea) so I wanted to try it too. That is the library that provides the functionality for the text-based user interface you encounter when using my CLI tool. I believe it is also the library that also facilitates the UI that tools built on the SST framework exhibit. 

A very cool library in my opinion. The LearnLM 1.5 Pro experimental model was crucial in my efforts to understand how to use this library. 

