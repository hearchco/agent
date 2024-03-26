<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![AGPL License][license-shield]][license-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/hearchco/hearchco">
    <img src="images/logo.svg" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">Hearchco</h3>

  <p align="center">
    Hearchco (pronounced /hɝːtʃko/, which stands for the Serbian word Hrčko) is a distributed and fast metasearch engine that respects your privacy. It is aimed to be a modern alternative to SearXNG and other metasearch engines. Hearchco collects results from multiple sources and deduplicates them while also ranking them depending on where they came from.
    <br />
    <!-- <a href="https://github.com/hearchco/hearchco"><strong>Explore the docs »</strong></a> -->
    <!-- <br /> -->
    <br />
    <a href="https://hearch.co">Hearch something!</a>
    ·
    <a href="https://github.com/hearchco/hearchco/issues">Report Bug</a>
    ·
    <a href="https://github.com/hearchco/hearchco/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <!-- <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul> -->
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <!-- <li><a href="#usage">Usage</a></li> -->
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

[![Product Name Screen Shot][product-screenshot]](https://hearch.co)

Hearchco's objectives:
  + speed - written in Go, built with caching and concurrency
  + privacy - no ads and no tracking, being a proxy to protect you from Big Tech
  + customizability - ranking, filtering and microserviced for a modular design
  + bangs - engine categories with specialized ranking for specialized searches
  + cuteness - lil hamster

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ### Built With -->

<!-- * [![Next][Next.js]][Next-url]
* [![Go][Go]][Go-url] -->

<!-- <p align="right">(<a href="#readme-top">back to top</a>)</p> -->



<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple example steps.

### Prerequisites

* go >= 1.22 - [how to install](https://go.dev/doc/install)

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/hearchco/hearchco
   ```
2. Install dependencies and generate code
   ```sh
   make install
   ```
3. Copy example config file
   ```sh
   cp hearchco_example.yaml hearchco.yaml
   ```
4. Start the router
   ```sh
   make run
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
<!-- ## Usage -->

<!-- Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources. -->

<!-- _For more examples, please refer to the [Documentation](https://example.com)_ -->

<!-- <p align="right">(<a href="#readme-top">back to top</a>)</p> -->



<!-- ROADMAP -->
## Roadmap

- [x] Basic search
- [x] Image search
- [ ] Search suggestions
- [ ] Search snippets

See the [open issues](https://github.com/hearchco/hearchco/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the AGPL-3.0 License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Aleksa Siriški - [@aleksasiriski](https://github.com/aleksasiriski)
? ? - [@k4lizen](https://github.com/k4lizen)
Matija Kljajić - [@matijakljajic](https://github.com/matijakljajic)

Project Link: [https://github.com/hearchco/hearchco](https://github.com/hearchco/hearchco)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [SearXNG](https://github.com/searxng/searxng) - for the inspiration and excellent reverse engineering
* [Colly](https://github.com/gocolly/colly) - for fast and reliable scraper
* [Rocketlaunchr](https://github.com/rocketlaunchr/google-search) - for showing how to scrape Google in Go

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/hearchco/hearchco.svg?style=for-the-badge
[contributors-url]: https://github.com/hearchco/hearchco/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/hearchco/hearchco.svg?style=for-the-badge
[forks-url]: https://github.com/hearchco/hearchco/network/members
[stars-shield]: https://img.shields.io/github/stars/hearchco/hearchco.svg?style=for-the-badge
[stars-url]: https://github.com/hearchco/hearchco/stargazers
[issues-shield]: https://img.shields.io/github/issues/hearchco/hearchco.svg?style=for-the-badge
[issues-url]: https://github.com/hearchco/hearchco/issues
[license-shield]: https://img.shields.io/github/license/hearchco/hearchco.svg?style=for-the-badge
[license-url]: https://github.com/hearchco/hearchco/blob/main/LICENSE.txt
[product-screenshot]: images/screenshot.png