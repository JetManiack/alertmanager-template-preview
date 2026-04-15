const alertmanagerFuncs = [
  "toUpper", "toLower", "title", "trimSpace", "join", "match", "safeHtml", 
  "safeUrl", "urlUnescape", "reReplaceAll", "stringSlice", "date", "tz", 
  "since", "humanizeDuration", "toJson", "list", "append", "dict"
].map(name => ({ label: name, type: "function" }));

const prometheusFuncs = [
  "toUpper", "toLower", "title", "trimSpace", "join", "match", "reReplaceAll",
  "humanize", "humanize1024", "humanizeDuration", "humanizeTimestamp", "humanizePercentage",
  "query", "first", "last", "value", "label"
].map(name => ({ label: name, type: "function" }));

const alertmanagerFields = [
  "Receiver", "Status", "Alerts", "NotificationReason", 
  "GroupLabels", "CommonLabels", "CommonAnnotations", "ExternalURL"
].map(name => ({ label: name, type: "variable" }));

const alertFields = [
  "Status", "Labels", "Annotations", "StartsAt", "EndsAt", 
  "GeneratorURL", "Fingerprint"
].map(name => ({ label: name, type: "variable" }));

const alertsMethods = [
  "Firing", "Resolved"
].map(name => ({ label: name, type: "method" }));

const kvMethods = [
  "SortedPairs", "Names", "Values"
].map(name => ({ label: name, type: "method" }));

const prometheusFields = [
  "Labels", "ExternalLabels", "ExternalURL", "Value", "Queries"
].map(name => ({ label: name, type: "variable" }));

// All possible variable completions (without leading dot)
const amAllFieldCompletions = [...alertmanagerFields, ...alertFields, ...alertsMethods, ...kvMethods];
const promAllFieldCompletions = [...prometheusFields];

// All possible variable completions with a leading dot for "no-dot" contexts
const amAllDottedCompletions = amAllFieldCompletions.map(c => ({
  ...c,
  label: "." + c.label,
  filterText: c.label 
}));

const promAllDottedCompletions = promAllFieldCompletions.map(c => ({
  ...c,
  label: "." + c.label,
  filterText: c.label 
}));

export function createTemplateCompletionSource(alertData, mode = 'alertmanager') {
  const isProm = mode === 'prometheus';
  const funcs = isProm ? prometheusFuncs : alertmanagerFuncs;
  const allFieldCompletions = isProm ? promAllFieldCompletions : amAllFieldCompletions;
  const allDottedCompletions = isProm ? promAllDottedCompletions : amAllDottedCompletions;

  return (context) => {
    // Match a word before the cursor (including a dot)
    let word = context.matchBefore(/\.?\w*/);
    
    if (!word || (word.from === word.to && !context.explicit)) {
      return null;
    }

    // Check if the word starts with a dot
    const startsWithDot = word.text.startsWith(".");
    
    if (startsWithDot) {
      // Possible nested path: .Field.SubField
      const beforeWord = context.state.sliceDoc(0, word.from);
      const parentMatch = beforeWord.match(/\.(\w+)$/);
      
      if (parentMatch) {
        const field = parentMatch[1];
        let source = null;
        
        if (isProm) {
          if (field === "Labels" && alertData?.labels) source = alertData.labels;
          else if (field === "ExternalLabels" && alertData?.externalLabels) source = alertData.externalLabels;
          else if (field === "Queries" && alertData?.queries) source = alertData.queries;
        } else {
          if (field === "CommonLabels" && alertData?.commonLabels) source = alertData.commonLabels;
          else if (field === "GroupLabels" && alertData?.groupLabels) source = alertData.groupLabels;
          else if (field === "CommonAnnotations" && alertData?.commonAnnotations) source = alertData.commonAnnotations;
        }
        
        if (source) {
          const keys = Object.keys(source).map(key => ({
            label: key,
            type: "property",
          }));
          return {
            from: word.from + 1, // Skip the dot
            to: context.pos,
            options: keys,
            validFor: /^\w*$/
          };
        }
      }

      // Default fields after a dot
      return {
        from: word.from + 1, // Skip the dot
        to: context.pos,
        options: allFieldCompletions,
        validFor: /^\w*$/
      };
    }

    // No dot -> Suggest functions AND top-level fields (starting with dot)
    return {
      from: word.from,
      to: context.pos,
      options: [...funcs, ...allDottedCompletions],
      validFor: /^\w*$/
    };
  };
}
