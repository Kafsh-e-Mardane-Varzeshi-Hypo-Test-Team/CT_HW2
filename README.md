# CT_HW2
Computer Technology Course - Homework 2

# [KeMV Online Judge Website](https://github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2)

This project is an online code judging platform, inspired by websites like Codeforces and Quera, where users can view questions, submit their solutions, and see the results of their submissions. The system provides functionalities such as user authentication, code submission handling, and an admin panel for managing questions.

## Features

### User Authentication & Roles

- Users can register and log in.
- Passwords are securely stored in the database using bcrypt.
- There are two types of users: regular users and admin users. Admins can manage questions, change user roles, and perform other privileged actions.
- Proper role-based access control is implemented to ensure that only admins can see and use admin functionalities.

### User Profiles

- Each user has a profile showing their activity, including statistics about their submissions.
- Admin users can change regular users’ roles (promote to admin or demote).
  
### Question Listing

- Questions are listed by their publish date (newest first).
- Pagination is implemented for question listing, with 10 questions displayed per page.

### Question Management

- Admin users can create, edit, and publish questions.
- Questions have the following attributes: title, content (statement), time and memory limits, input/output data.
- Questions are initially in a draft state and can only be published by an admin.

### Submissions

- Users can submit solutions (written in Golang) to questions.
- Submissions go through a judging process, where the code is tested against predefined cases.
- The result of a submission can be one of the following: Correct (OK), Compile Error, Wrong Answer, Memory Limit Exceeded, Time Limit Exceeded, or Runtime Error.

### Code Runner Service

- A dedicated "code-runner" service handles the execution of user-submitted code in isolated Docker containers.
- This service ensures security by limiting resource usage (CPU, memory) and disallowing network access for submitted code.
  
### Docker Deployment

- The project uses Docker to containerize the application, including the code runner service.
- Docker Compose is used to set up an internal network that connects all necessary services: the main app, the database, and the code runner service.

### Load Testing & Performance

- A script for load testing is used to simulate high traffic and ensure that the system handles large amounts of data and submissions efficiently.
- You can see load testing results from [here](https://github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/blob/feature/load-test/load_test_report.md)!

## Project Phases

The whole project was dissected into 7 phases:

1. **Frontend** - [Mohammad Barkatin](https://github.com/mammedbrk) implemented the frontend.
2. **Authentication** - [Mehdi Taheri](https://github.com/Mefi22) handled the user authentication and session management.
3. **Database Design** - [Parisa Sadat Mousavi](https://github.com/parisam83) was responsible for designing the database schema.
4. **Implementing Site Services** - [Mehdi Taheri](https://github.com/Mefi22) worked on implementing the core services of the site.
5. **Code Runner Service & Integration** - [Mohammad Barkatin](https://github.com/mammedbrk) and [Mehdi Taheri](https://github.com/Mefi22) developed and integrated the code runner service.
6. **Docker Deployment** - This task remains to be completed.
7. **Database Population & Load Testing** - [Parisa Sadat Mousavi](https://github.com/parisam83) populated the database and conducted load testing.

