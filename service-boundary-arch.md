## How to build a cloud native product

There are social barriers you need to overcome to move from being a distributed monolith to a cloud native product.

### Legacy systems

- Build an adaptor team that can hide technical debt from old systems as an API
- The adaptor team allows new services to be built without forklifts

### Security

- All servieces should be runnable insecurely
- All services should also be runnable in an extremely secure way
- If you dont agree with the first point: dedicate a team to building developer tooling, b/c
otherwise your slowing people down

### Expertise

Breaking teams into automonomous subunits doesnt mean you dont have expertise informing your services.  
Teams should regularly get help from experts on other teams.  
Make sure each team has at least one person whose social and willing to collaborate.

### Languages

- agree on one language that all teams agree to communicate with.  Python or Go or Java.
- this language is only for communication and idea sharing.  teams can implement services in any other number of languages.

### Architecture and requirements

- you're architecture is containers over REST, everything else is constantly changing.  Don't freeze it.
- you're requirements are API specifications.
- your choice of database is not an architecture, neither is your language or UML diagram.
- Your architecture is what your app does, not what libraries or databases or tools you use.
- the team builidng the service builds the first specification and a demo app.

### Performance testing

- performance bugs are filed as scripts or programs that run in automation

### Documentation

No documentation other then swagger is explicitly required for any service. 
Teams should however be evaluated based on what docs/videos/demo apps they voluntarily 
provide to others consuming their service.

- Teams market their services by sharing examples and swagger docs.
- teams consuming the above service do so by startnig w/ the demo app as a template.

### Ownership

Ownership is everything
- teams innovate and build services
- teams build better services when they own them
- teams build crappy services when someone else owns them
- as teams succeed at one service, let them take on new services

### Keep it flat
- no team exceeds 4 people
- teams can be as small as one person
- if a team is overloaded, find a way to split their service up, dont stack the decks
- Integration and QA is a team in and of itself.
- Salary isn't coupled to how many people you order, but how experienced you are.
  - Individual contributors with domain expertise make the most 
  
### Integration team rotation
- One team  integrations the services and runs them in automation
- This team has a rotating seat that has a member from other teams

### Product management runs the product 
- PM is QA, and understands the product
- The integration team files bugs against microservices as needed.

### Promotion
- Engineers on teams can ask, as they become more senior, to have their own team.  
- That conversation should be data driven.

