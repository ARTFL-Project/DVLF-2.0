(function() {
    "use strict";

    angular
        .module("DVLF")
        .directive('apropos', apropos);


    function apropos($rootScope) {
        return {
            templateUrl: 'app/components/apropos/apropos.html',
            replace: true
        };
    }
})();
