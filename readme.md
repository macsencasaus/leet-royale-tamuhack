# leet-royale-tamuhack
The following is our original devpost submission for **TAMUhack 2025**, it can also be found on [Devpost](https://devpost.com/software/leet-royale).

---

## Inspiration

Every computer science major knows of the dreaded _whiteboard interview_. Preparing is important, but often difficult and boring. Using **Leet Royale**, studying becomes just as fun as achieving a Victory Royale. 

## What it does

By connecting with friends and strangers alike, Leet Royale pushes your skills as a programmer and prepares you for your next interview. Go through 4 rigorous rounds of coding questions where only the best will make it to the next round. By the final round, only a single person will be left standing!

## How we built it

We built this by splitting it into 4 major parts: UI, compilation, problem set design, and networking. The UI frontend was created using React and the backend including the problem set design, compilation, and networking were all in golang. The problem set design automatically generates templates for the target languages supported by the website, the compilation utilizes clang, node, and python to observe output and errors,  while the networking used websockets to create peer-to-server connections with the server to communicate events from individuals to all users.

## Challenges we ran into

We expected this to be challenging, but the project presented major issues out of our predictions. A particularly interesting challenge that we faced was the need to avoid writing to disk. Since users may upload programs for evaluation very quickly, it would put unnecessary strain on us to write to disk (or even temp files) for every uploaded program. Since we already have to compile and run them, it was most efficient to process them using only memory buffers controlled by the web server. To accomplish this, we used low-level LLVM tools to compile and run chunks of LLVM bitcode without producing a complete ELF file.

## Accomplishments that we're proud of

We're particularly proud of the robust messaging system that can support over 40 users on websockets. Behind the scenes, there's nearly 16 unique types of messages that have to be handled by different parts of the system including joining lobbies, getting test case results, eliminating users, and viewing other players' statuses.  

## What we learned

In this project, we learned that frequent and precise communication about specifications and interfaces is extremely important for allowing each team member to succeed. By developing a clear understanding of how each part of the code should communicate with other parts, we were able to split up the project into small parts that each team member could do separately and quickly, while still producing a well-integrated full-project.

## What's next for Leet Royale

The next up for Leet Royale is definitely a shop. We imagine a wholly unique system where users can buy buffs for themself (like assistance from ChatGPT or extra time) or debuffs for opponents (like making a bouncing DvD logo that blocks their view, forcing their editor into light mode, forcing them to move exclusively with the arrow keys, and removing semicolons from their code). We believe this would be a fun addition to offset the rather skill based game. A chat system would also be a fun addition and could lead to some interesting interactions.
