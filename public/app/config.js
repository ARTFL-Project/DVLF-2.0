(function() {
    "use strict";

    angular
        .module("DVLF")
        .config(DVLFConfig);

    function DVLFConfig($routeProvider, $locationProvider) {
        $routeProvider.
        when('/mot/:queryTerm', {
            templateUrl: 'app/components/results/results.html',
            controller: 'ResultsController',
            controllerAs: 'Results'
        }).
        when('/apropos', {
            templateUrl: 'app/components/apropos/apropos.html'
        }).
        when('/definition', {
            templateUrl: 'app/components/newDefinition/newDefinition.html',
            controller: "NewDefinitionController",
            controllerAs: "NewDefinition"
        }).
        otherwise({
            redirectTo: '/'
        });
        $locationProvider.html5Mode({
            enabled: true
        });
    }
})()
