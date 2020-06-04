# CHANGELOG
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html)

<a name="unreleased"></a>
## [Unreleased]
### Bug Fixes
- Fix white company page if it had a relationship with a deleted company

### Code Refactoring
- Move from `mux` and `negroni` to `gin`

### Features
- Implement authentication
- Implement authorization
- Implement user registration
- Implement user and permission management UI
- Implement the changelog for companies
- Improve user experience
- Drop concept of company sector
- Create concept of company verticals
- Add the suggestion of the relationship tier hovering on the icon
- Allow asymmetrical relationships
- Implement uniqueness guarantees for companies and relationships
- Add the field "notes" to relationships

<a name="0.1.0"></a>
## [0.1.0] - 2020-05-23
### Features
- Initial version
