(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('examples', examples);

    function examples($http, $timeout, $log, $rootScope) {
        return {
            templateUrl: "app/components/examples/examples.html",
            link: function(scope, element, attrs) {
            
            }
        }
    }
})();
