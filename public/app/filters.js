(function() {
    "use strict";

    angular
        .module("DVLF")
        .filter("isEmpty", isEmpty)
        .filter("unsafe", unSafe)
        .filter('encodeURIComponent', encode)
        .filter("sortByVote", sortByVote);


    function isEmpty() {
        return function(obj) {
            if (angular.element.isEmptyObject(obj)) {
                return false;
            } else {
                return true;
            }
        };
    }

    function unSafe($sce) {
        return $sce.trustAsHtml;
    }

    function encode($window) {
        return $window.encodeURIComponent;
    }

    function sortByVote() {
        return function(examples) {
            var orderedExamples = [];
            var negativeExamples = [];
            var otherExamples = [];
            for (var i=0; i < examples.length; i+=1) {
                if (examples[i].score > 0) {
                    orderedExamples.push(examples[i]);
                } else if (examples[i].score < 0){
                    negativeExamples.push(examples[i])
                } else {
                    otherExamples.push(examples[i]);
                }
            }
            orderedExamples.sort(function(a, b) {
                return b.score - a.score;
            });
            negativeExamples.sort(function(a, b) {
                return b.score - a.score;
            });

            var allExamples = orderedExamples.concat(otherExamples, negativeExamples);
            return allExamples;
        }
    }
})();
