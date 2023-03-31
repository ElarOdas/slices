# Slices

## What is Slices

The slices package implements functions to handle slices that are commonly found in functional languages
Including:

-   map
-   filter
-   reduce
-   some
-   every
-   flat

Of note: The [official Slices](golang.org/x/exp/slices) project provides other common slice functions like search,index, replace and sort.

## How does it work

To realize the above mentioned functions the (newly) added generics were used. As this is a go package concurrency is also
used where applicable.
It is not possible to directly assign methods to slices and go reservers some keywords like map.
For now we circumvent this problem by using a (Method)**Slice** naming convention.

## Usage

For now the intended user is only me. Purpose of this package is to reduce the boilerplate while working with slices in go in my other projects.
Secondary purpose for this project is to practice the usage of generics and concurrency.

## How to

### Develop

Feel free to fork the project! For now I can not guarantee further work on this project

### Use

This package is currently highly experimental and should not be used in a production environment.

### Test

Right now no testing is implemented, but I intend to add testing at a later point.

## Plans

-   more extensive documentation
-   example project (see tests for now)
-   better error handling
    Right now I am busy with my [master thesis]().
